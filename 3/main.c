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

struct part {
  // inclusive
  int istart;
  // exclusive
  int iend;
  int value;
};

struct part parse_part_number(char *line, int starti) {
  int consecutive = 0;
  char number[4096];
  int i;
  for (i = starti; i < strlen(line); i++) {
    char c = line[i];
    if (c >= '0' && c <= '9') {
      // printf("got digit: %c\n", c);
      number[consecutive] = c;
      consecutive++;
    } else {
      number[consecutive] = '\0';
      if (consecutive > 0) {
        // printf("closed number: %s\n", number);
        break;
      }
    }
  }

  if (consecutive == 0) {
    struct part p = {-1, -1, -1};
    return p;
  }
  struct part p = {i - strlen(number), i, atoi(number)};
  return p;
}

int is_symbol(char c) {
  if (c == '!')
    return 1;
  if (c == '@')
    return 1;
  if (c == '#')
    return 1;
  if (c == '$')
    return 1;
  if (c == '%')
    return 1;
  if (c == '^')
    return 1;
  if (c == '&')
    return 1;
  if (c == '*')
    return 1;
  if (c == '-')
    return 1;
  if (c == '_')
    return 1;
  if (c == '=')
    return 1;
  if (c == '+')
    return 1;
  if (c == '?')
    return 1;
  if (c == '>')
    return 1;
  if (c == '<')
    return 1;
  if (c == '/')
    return 1;
  return 0;
}

int is_part_number(struct part p, char *pline, char *cline, char *nline) {
  int xmin = max(0, p.istart - 1);
  int xmax = min(strlen(nline), p.iend + 1);

  for (int i = xmin; i < xmax; i++) {
    char c = pline[i];
    // printf("checking symbol %c for part %d\n", c, p.value);
    if (is_symbol(c) > 0) {
      return 1;
    }
  }

  char c = cline[xmin];
  // printf("checking symbol %c for part %d\n", c, p.value);
  if (is_symbol(c) > 0) {
    return 1;
  }

  c = cline[xmax - 1];
  // printf("checking symbol %c for part %d\n", c, p.value);
  if (is_symbol(c) > 0) {
    return 1;
  }

  for (int i = xmin; i < xmax; i++) {
    char c = nline[i];
    // printf("checking symbol %c for part %d\n", c, p.value);
    if (is_symbol(c) > 0) {
      return 1;
    }
  }

  return 0;
}

int parse_line(char *l1, char *l2, char *l3) {
  printf("1: %s", l1);
  printf("2: %s", l2);
  printf("3: %s", l3);
  printf("\n");

  int line_parts_sum = 0;

  struct part p = {-1, -1, -1};
  int starti = 0;
  do {
    p = parse_part_number(l2, starti);
    printf("number start: %d, end: %d, value: %d\n", p.istart, p.iend, p.value);
    if (p.iend > -1) {
      starti = p.iend;

      if (is_part_number(p, l1, l2, l3) > 0) {
        printf("%d is part number\n", p.value);
        line_parts_sum += p.value;
      } else {
        printf("%d is NOT part number\n", p.value);
      }
    }
  } while (p.value != -1);

  return line_parts_sum;
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
  while ((read = getline(&line, &len, fp))) {
    i++;
    printf("\nline %d, read: %ld \n", i, read);

    strcpy(lines[0], lines[1]);
    strcpy(lines[1], lines[2]);
    strcpy(lines[2], line);

    if (i < 1) {
      continue;
    } else if (i == 1) {
      int l = read;
      for (int i = 0; i < read; i++) {
        lines[0][i] = '.';
      }
      lines[0][l - 1] = '\n';
      lines[0][l] = '\0';
    } else if (read == -1) {
      int l = strlen(lines[1]);
      for (int i = 0; i < l; i++) {
        lines[2][i] = '.';
      }
      lines[2][l - 1] = '\n';
      lines[2][l] = '\0';
    }

    sum += parse_line(lines[0], lines[1], lines[2]);

    if (read == -1)
      break;
  }

  printf("sum of parts is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
