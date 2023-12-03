#include <ctype.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifndef max
#define max(a, b) (((a) > (b)) ? (a) : (b))
#endif

#ifndef min
#define min(a, b) (((a) < (b)) ? (a) : (b))
#endif

void parse_line(char *l1, char *l2, char *l3) {
  printf("1: %s", l1);
  printf("2: %s", l2);
  printf("3: %s", l3);
  printf("\n");

  int consecutive = 0;
  char number[4096];
  for (int i = 0; i < strlen(l1); i++) {
    char c = l1[i];
    if (c >= '1' && c <= '9') {
      printf("got digit: %c\n", c);
      number[consecutive] = c;
      consecutive++;
    } else {
      number[consecutive] = '\0';
      if (consecutive > 0)
        printf("closed number: %s\n", number);
      consecutive = 0;
    }
  }
}

int main(void) {
  char *line = NULL;
  size_t len = 0;
  ssize_t read;

  FILE *fp = fopen("input.txt", "r");
  if (fp == NULL)
    exit(EXIT_FAILURE);

  unsigned sum = 0;
  char lines[3][4096];
  int i = -1;
  while ((read = getline(&line, &len, fp)) != -1) {
    i++;

    strcpy(lines[0], lines[1]);
    strcpy(lines[1], lines[2]);
    strcpy(lines[2], line);
    if (i > 1) {
      parse_line(lines[0], lines[1], lines[2]);
    }
  }

  printf("sum of ids is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
