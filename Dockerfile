FROM scratch
COPY build/overtime_linux /overtime
USER 7000
ENTRYPOINT ["/overtime"]