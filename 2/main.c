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

const int rgb[3] = {12, 13, 14};

char *ltrim(char *s) {
  while (isspace(*s))
    s++;
  return s;
}

char *rtrim(char *s) {
  char *back = s + strlen(s);
  while (isspace(*--back))
    ;
  *(back + 1) = '\0';
  return s;
}

char *trim(char *s) { return rtrim(ltrim(s)); }

// true if ok
int process_color(char *color) {
  printf("color: %s\n", color);

  if (strstr(color, "red")) {
    int v;
    sscanf(color, "%d", &v);
    printf("got red %d\n", v);
    if (v > rgb[0])
      return 0;
  } else if (strstr(color, "green")) {
    int v;
    sscanf(color, "%d", &v);
    printf("got green %d\n", v);
    if (v > rgb[1])
      return 0;
  } else if (strstr(color, "blue")) {
    int v;
    sscanf(color, "%d", &v);
    printf("got blue %d\n", v);
    if (v > rgb[2])
      return 0;
  }
  return 1;
}

// true if ok
int process_token(char *token) {
  printf("token: %s\n", token);

  char *end_str;
  char *color = strtok_r(token, ",", &end_str);
  while (color != NULL) {
    int ok = process_color(trim(color));
    if (ok == 0)
      return 0;
    color = strtok_r(NULL, ",", &end_str);
  }
  return 1;
}

int parse_line(char *line) {
  printf("parsing: %s", line);
  line = strstr(line, ":") + 2;
  printf("skipped: %s", line);

  char *end_str;
  char *token = strtok_r(line, ";", &end_str);
  while (token != NULL) {
    int ok = process_token(trim(token));
    if (ok == 0)
      return 0;
    token = strtok_r(NULL, ";", &end_str);
  }

  return 1;
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
    int ok = parse_line(line);
    if (ok > 0) {
      printf("id %d was ok\n\n", id);
      sum += id;
    } else {
      printf("id %d was not ok\n\n", id);
    }
  }

  printf("sum of ids is: %d\n", sum);

  fclose(fp);
  if (line)
    free(line);
  exit(EXIT_SUCCESS);
}
