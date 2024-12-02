#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <stdint.h>

#define buffer_size 8192

int main(void) {
    uint64_t buf[buffer_size >> 3];
    uint32_t acc;
    size_t i;

    int const fd = open("nya", 0, "r");
    acc = read(fd, buf, buffer_size);
    i = acc >> 3;

    const uint64_t mask = 0x0101010101010101;

    do {
        i--;
        // Find all ) in the current word by looking for the trailing bit in )
        uint64_t current_block = buf[i] & mask;
        // Count how many there are by folding and adding the bit array
        current_block += current_block >> 32;
        current_block += current_block >> 16;
        current_block += current_block >>  8;
        // Subtract each one _twice_. We assumed at the start everything was (,
        // so we have to remove the ( by subtracting one, and then add the )
        // by also subtracting one.
        acc -= (current_block & 15) << 1;
    } while (i != 0);

    printf("%d\n", acc);
}
