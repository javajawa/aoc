#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>

#define buffer_size 8192

int main(void) {
    int fd;
    int len;
    unsigned char buf[buffer_size];

    fd = open("nya", 0, "r");
    len = read(fd, buf, buffer_size);

    int stepsTaken = 0;
    int downStepsTaken = 0;

    while (stepsTaken < len) {
        downStepsTaken += (buf[stepsTaken] & 1);
        stepsTaken += 1;
        if (downStepsTaken << 1 > stepsTaken) {
            printf("%d\n", stepsTaken);
            return 0;
        }
    }

    printf("Basement never reached\n");
    return 1;
}
