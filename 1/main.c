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

unsigned concatenate(unsigned x, unsigned y) {
  unsigned pow = 10;
  while (y >= pow)
    pow *= 10;
  return x * pow + y;
}

struct numbers {
  int first;
  int last;
};

struct numbers get_indices_for_str(char *line, char *needle) {
  int iname_min = INT_MAX;
  int iname_max = INT_MIN;
  char *currline = line;
  char *pname = currline;
  while (pname != NULL) {
    printf("searching for %s in %s", needle, currline);
    pname = strstr(currline, needle);
    if (pname != NULL) {
      printf("got %s", pname);
      int pos = pname - line;
      if (pos > iname_max)
        iname_max = pos;
      if (pos < iname_min)
        iname_min = pos;
      currline = pname + strlen(needle);
    }
  }

  struct numbers res = {iname_min, iname_max};
  return res;
}

struct numbers get_indices_for_num(char *line, int num) {
  char names[][6] = {"one", "two",   "three", "four", "five",
                     "six", "seven", "eight", "nine"};
  char *name = names[num - 1];
  struct numbers pname = get_indices_for_str(line, name);
  char number[2];
  sprintf(number, "%d", num);
  struct numbers pnumber = get_indices_for_str(line, number);

  struct numbers res = {min(pname.first, pnumber.first),
                        max(pname.last, pnumber.last)};
  printf("positions for %s are %d and %d\n", name, res.first, res.last);
  return res;
}

struct numbers parse_line(char *line) {
  printf("parsing: %s", line);

  int imin = INT_MAX;
  int imax = INT_MIN;
  int nummin;
  int nummax;
  for (int i = 1; i < 10; i++) {
    struct numbers inds = get_indices_for_num(line, i);
    if (inds.first < imin) {
      imin = inds.first;
      nummin = i;
    }
    if (inds.last > imax) {
      imax = inds.last;
      nummax = i;
    }
  }

  struct numbers res = {nummin, nummax};
  return res;
}

int main(void) {
  char *line = NULL;
  size_t len = 0;
  ssize_t read;

  FILE *fp = fopen("input.txt", "r");
  if (fp == NULL)
    exit(EXIT_FAILURE);

  unsigned sum = 0;
  while ((read = getline(&line, &len, fp)) != -1) {

    struct numbers num = parse_line(line);
    printf("first: %d last: %d\n", num.first, num.last);
    unsigned code = concatenate(num.first, num.last);
    printf("%d\n", code);
    sum += code;
  }

  printf("sum of codes is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
