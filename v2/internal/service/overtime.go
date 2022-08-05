package service

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/api/v2/pkg/utils"
)

type MainService struct {
	db        *data.Db
	instances map[uint]pkg.OvertimeService
}

type Service struct {
	user     *pkg.User
	db       *data.Db
	readonly bool
}

func Init(db *data.Db) *MainService {
	return &MainService{
		db:        db,
		instances: map[uint]pkg.OvertimeService{},
	}
}

func (s *MainService) GetOrCreateInstanceForUser(user *pkg.User) pkg.OvertimeService {
	if ser, exists := s.instances[user.ID]; exists {
		return ser
	}

	userInstance := Service{
		user: user,
		db:   s.db,
	}

	s.instances[user.ID] = &userInstance

	return &userInstance
}

func (s *MainService) GetOrCreateReadonlyInstanceForUser(user *pkg.User) pkg.OvertimeService {
	return &Service{
		user:     user,
		db:       s.db,
		readonly: true,
	}
}

func (s *Service) calcOvertimeAndActivetime(start time.Time, end time.Time) (int64, int64, error) {
	overtimeInMinutes := int64(0)
	activeTimeInMinutes := int64(0)
	now := time.Now()

	st := start
	for {
		if st.Unix() >= end.Unix() {
			break
		}
		be := utils.DayStart(st)
		en := utils.DayEnd(st)
		if end.Unix() < en.Unix() {
			en = end
		}
		isNowDay := (be.Year() == now.Year() && be.Month() == now.Month() && be.Day() == now.Day())
		if !isNowDay {
			wd, err := s.db.GetWorkDay(be, s.user.ID)
			if err != nil {
				log.Info(err)
			}
			if wd != nil && err == nil {
				activeTimeInMinutes += wd.ActiveTime
				overtimeInMinutes += wd.Overtime
				st = st.AddDate(0, 0, 1)
				continue
			}
		}

		dayWorkTimeInMinutes, err := s.CalcDailyWorktime(be)
		if err != nil {
			log.Debug(err)
			return 0, 0, err
		}

		at, err := s.SumActivityBetweenStartAndEndInMinutes(be, en)
		if err != nil {
			log.Debug(err)
			return 0, 0, err
		}

		ft, isLegal, err := s.SumHolidaysBetweenStartAndEndInMinutes(be, en)
		if err != nil {
			log.Debug(err)
			return 0, 0, err
		}
		dayOvertimeInMinutes := at + ft - int64(dayWorkTimeInMinutes)
		if ft > 0 && dayWorkTimeInMinutes == 0 && !isLegal {
			dayOvertimeInMinutes = at
		}
		if !isNowDay {
			err = s.db.SaveWorkDay(
				&data.WorkDayDB{
					WorkDay: pkg.WorkDay{
						InputWorkDay: pkg.InputWorkDay{
							Day:        be,
							Overtime:   dayOvertimeInMinutes,
							ActiveTime: at,
							UserID:     s.user.ID,
						},
						IsHoliday: ft > 0,
					},
				})
			if err != nil {
				return 0, 0, err
			}
		}
		overtimeInMinutes += dayOvertimeInMinutes
		activeTimeInMinutes += at
		st = st.AddDate(0, 0, 1)
	}

	return activeTimeInMinutes, overtimeInMinutes, nil
}

func (s *Service) CalcOverview(day time.Time) (*pkg.Overview, error) {
	yyyy, mm, dd := day.Date()
	wd := day.Weekday()
	wdNumber := weekDayToInt(wd)
	// This year
	yStart := time.Date(yyyy, 01, 01, 0, 0, 0, 0, day.Location())
	yat, yot, err := s.calcOvertimeAndActivetime(yStart, day)
	if err != nil {
		return nil, err
	}

	holidays, err := s.CountUsedHolidaysBetweenStartAndEnd(yStart, day)
	if err != nil {
		return nil, err
	}

	// This month
	mStart := time.Date(yyyy, mm, 01, 0, 0, 0, 0, day.Location())
	mat, mot, err := s.calcOvertimeAndActivetime(mStart, day)
	if err != nil {
		return nil, err
	}
	// This week
	wStart := time.Date(yyyy, mm, dd-wdNumber+1, 0, 0, 0, 0, day.Location())
	wat, wot, err := s.calcOvertimeAndActivetime(wStart, day)
	if err != nil {
		return nil, err
	}
	// This day
	dStart := time.Date(yyyy, mm, dd, 0, 0, 0, 0, day.Location())
	at, ot, err := s.calcOvertimeAndActivetime(dStart, day)
	if err != nil {
		return nil, err
	}

	_, wn := day.ISOWeek()
	o := &pkg.Overview{
		Date:                         day,
		WeekNumber:                   wn,
		ActiveTimeThisDayInMinutes:   at,
		ActiveTimeThisWeekInMinutes:  wat,
		ActiveTimeThisMonthInMinutes: mat,
		ActiveTimeThisYearInMinutes:  yat,
		OvertimeThisDayInMinutes:     ot,
		OvertimeThisWeekInMinutes:    wot,
		OvertimeThisMonthInMinutes:   mot,
		OvertimeThisYearInMinutes:    yot,
		UsedHolidays:                 int(holidays),
		HolidaysStillAvailable:       int(s.user.NumHolidays - holidays),
	}
	cra, err := s.db.GetRunningActivityByUserID(s.user.ID)
	if err == nil {
		o.ActiveActivity = &cra.Activity
	}

	return o, nil
}
