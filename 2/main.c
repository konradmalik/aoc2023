#include "../lib/minmax.c"
#include "../lib/trim.c"
#include <ctype.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void process_color(char *color, int *rgb) {
  printf("color: %s\n", color);

  if (strstr(color, "red")) {
    int v;
    sscanf(color, "%d", &v);
    printf("got red %d\n", v);
    rgb[0] = max(rgb[0], v);
  } else if (strstr(color, "green")) {
    int v;
    sscanf(color, "%d", &v);
    printf("got green %d\n", v);
    rgb[1] = max(rgb[1], v);
  } else if (strstr(color, "blue")) {
    int v;
    sscanf(color, "%d", &v);
    printf("got blue %d\n", v);
    rgb[2] = max(rgb[2], v);
  }
}

void process_token(char *token, int *rgb) {
  printf("token: %s\n", token);

  char *end_str;
  char *color = strtok_r(token, ",", &end_str);
  while (color != NULL) {
    process_color(trim(color), rgb);
    color = strtok_r(NULL, ",", &end_str);
  }
}

void parse_line(char *line, int *rgb) {
  printf("parsing: %s", line);
  line = strstr(line, ":") + 2;
  printf("skipped: %s", line);

  char *end_str;
  char *token = strtok_r(line, ";", &end_str);
  while (token != NULL) {
    process_token(trim(token), rgb);
    token = strtok_r(NULL, ";", &end_str);
  }
}

int main(void) {
  char *line = NULL;
  size_t len = 0;
  ssize_t read;

  FILE *fp = fopen("input.txt", "r");
  if (fp == NULL)
    exit(EXIT_FAILURE);

  int id = 0;
  unsigned sum = 0;
  while ((read = getline(&line, &len, fp)) != -1) {
    id++;

    int rgb[3] = {0, 0, 0};
    parse_line(line, rgb);
    printf("line rgb: %d, %d, %d\n", rgb[0], rgb[1], rgb[2]);
    int power = rgb[0] * rgb[1] * rgb[2];
    printf("line power: %d\n", power);
    printf("\n");
    sum += power;
  }

  printf("sum of ids is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
