#include "../lib/trim.c"
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#define ELEMS 100
#define CARDS 500

void parse_card(char *card, int *result) {

  int i = 0;
  char *end_str;
  char *snum = strtok_r(card, " ", &end_str);
  while (snum != NULL) {
    int num = atoi(snum);
    result[i] = num;
    snum = strtok_r(NULL, " ", &end_str);
    i++;
  }
}

int parse_line(char *line) {
  int matches = 0;
  printf("parsing: %s", line);
  line = strstr(line, ":") + 1;

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

    for (int j = 0; j < ELEMS; j++) {
      int wnum = winnum[j];
      if (wnum == -1)
        break;

      if (wnum == anum) {
        matches++;
      }
    }
  }

  return matches;
}

int main(void) {
  char *line = NULL;
  size_t len = 0;
  ssize_t read;

  FILE *fp = fopen("input.txt", "r");
  if (fp == NULL)
    exit(EXIT_FAILURE);

  int line_x[CARDS];
  memset(line_x, 0, sizeof(line_x));

  int card = 0;
  unsigned sum = 0;
  while ((read = getline(&line, &len, fp)) != -1) {
    line_x[card] += 1;
    int matches = parse_line(line);
    printf("line matches: %d\n", matches);
    printf("\n");

    for (int m = 1; m <= matches; m++) {
      line_x[card + m] += line_x[card];
    }
    card++;
  }

  for (int i = 0; i < card; i++) {
    printf("card %d has %d copies\n", i + 1, line_x[i]);
    sum += line_x[i];
  }

  printf("sum of cards is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
