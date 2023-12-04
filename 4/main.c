#include "../lib/trim.c"
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#define ELEMS 100

void parse_card(char *card, int *result) {

  int i = 0;
  char *end_str;
  char *snum = strtok_r(card, " ", &end_str);
  while (snum != NULL) {
    int num = atoi(snum);
    printf("num: %d\n", num);
    result[i] = num;
    snum = strtok_r(NULL, " ", &end_str);
    i++;
  }
}

int parse_line(char *line) {
  int line_score = 0;
  printf("parsing: %s", line);
  line = strstr(line, ":") + 1;
  printf("skipped: %s", line);

  char *end_str;
  char *winning = trim(strtok_r(line, "|", &end_str));
  printf("winning: %s\n", winning);
  char *actual = strtok_r(NULL, "|", &end_str);
  printf("actual: %s\n", actual);

  int winnum[ELEMS];
  memset(winnum, -1, sizeof(winnum));
  parse_card(winning, winnum);

  int actnum[ELEMS];
  memset(actnum, -1, sizeof(actnum));
  parse_card(actual, actnum);

  for (int i = 0; i < ELEMS; i++) {
    int anum = actnum[i];
    if (anum == -1)
      break;

    printf("compare %d\n", anum);
    for (int j = 0; j < ELEMS; j++) {
      int wnum = winnum[j];
      if (wnum == -1)
        break;

      if (wnum == anum) {
        if (line_score == 0)
          line_score = 1;
        else
          line_score *= 2;
      }
    }
  }

  return line_score;
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
    int line_score = parse_line(line);
    printf("line score: %d\n", line_score);
    printf("\n");
    sum += line_score;
  }

  printf("sum of scores is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
