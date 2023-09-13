// gcc -Wall -g -O -o thraff -lpthread ./thraff.c

#define _GNU_SOURCE
#include <pthread.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MAX_NUM_THREADS 128

static void set_thread_affinity(pthread_t t, int cpu) {
  cpu_set_t cpuset;

  CPU_ZERO(&cpuset);
  CPU_SET(cpu, &cpuset);
  pthread_setaffinity_np(t, sizeof(cpu_set_t), &cpuset);
}

static void print_thread_affinity(pthread_t t, int ID) {
  cpu_set_t cpuset;
  int i;

  pthread_getaffinity_np(t, sizeof(cpu_set_t), &cpuset);

  for (i = 0; i < CPU_SETSIZE; i++) {
    if (CPU_ISSET(i, &cpuset))
      printf("Thread %2d bound to CPU %d\n", ID, i);
  }
}

static void *worker_thread_func(void *args) {
  int ID = *(int *)args;

  print_thread_affinity(pthread_self(), ID);

  while (1) {
    ; // TODO: better busy loop
  }

  return NULL;
}

int atoilist(char **pieces, int *out, int outlen) {
  int val, i = 0;
  if (pieces == NULL) {
    return -1;
  }
  int outpos = 0;
  for (i = 0; pieces[i] != NULL; i++) {
    val = atoi(pieces[i]);
    if (val == -1) {
      return -1;
    }
    out[outpos++] = val;
    if (outpos >= outlen) {
      break;
    }
  }
  return outpos;
}

char **strsplit(const char *str, char sep, size_t *pieces_num) {
  const char *begin = str, *end = NULL;
  char **pieces = NULL, *pc = NULL;
  size_t i = 0, n = 2;
  int failed = 0;

  if (!str || !strlen(str)) {
    return NULL;
  }

  while (begin != NULL) {
    begin = strchr(begin, sep);
    if (begin != NULL) {
      begin++;
      n++;
    }
  }

  pieces = malloc(n * sizeof(char *));
  if (!pieces) {
    return NULL;
  }

  begin = str;
  while (begin != NULL) {
    size_t len;

    end = strchr(begin, sep);
    if (end != NULL) {
      len = (end - begin);
    } else {
      len = strlen(begin);
    }
    if (len > 0) {
      pc = strndup(begin, len);
      if (pc == NULL) {
        failed = 1;
        break;
      } else {
        pieces[i] = pc;
        i++;
      }
    }
    if (end != NULL) {
      begin = end + 1;
    } else {
      break;
    }
  }

  if (failed) {
    /* one or more copy of pieces failed */
    free(pieces);
    pieces = NULL;
  } else {                /* i == n - 1 -> all pieces copied */
    pieces[n - 1] = NULL; /* end marker */
    if (pieces_num != NULL) {
      *pieces_num = i;
    }
  }
  return pieces;
}

void strfreev(char **pieces) {
  int i = 0;
  if (pieces == NULL) {
    return;
  }
  for (i = 0; pieces[i] != NULL; i++) {
    free(pieces[i]);
  }
  free(pieces);
}


int main(int argc, char *argv[]) {
  int num_worker_threads = -1;
  int IDs[MAX_NUM_THREADS] = { 0 }, i;
  pthread_t threads[MAX_NUM_THREADS];
  int cpus[MAX_NUM_THREADS] = { 0 };
  size_t num_cpus = 0;
  char **cpu_list = NULL;
  int opt, ret;

  while ((opt = getopt(argc, argv, "c:w:h")) != -1) {
    switch (opt) {
    case 'c':
      cpu_list = strsplit(optarg, ',', &num_cpus);
      if (cpu_list == NULL) {
        fprintf(stderr, "invalid number of cpus: %s\n", optarg);
        exit(1);
      }
      ret = atoilist(cpu_list, cpus, (int)num_cpus);
      if (ret == -1) {
        fprintf(stderr, "error parsing cpus: %s\n", optarg);
        exit(1);
      }
      printf("using %d cpus for affinity\n", ret);
      strfreev(cpu_list);
      break;
    case 'w':
      num_worker_threads = atoi(optarg);
      if (num_worker_threads == -1) {
        fprintf(stderr, "invalid number of cpus: %s\n", optarg);
        exit(1);
      }
      break;
    case 'h':
      fprintf(stderr, "usage: %s [-c][-w][-h]\n", argv[0]);
      exit(0);
    }
  }

  if (num_worker_threads == -1) {
    fprintf(stderr, "missing worker thread count\n");
    exit(1);
  }
  printf("using %zd cores, %d worker threads\n", num_cpus, num_worker_threads);

  for (i = 0; i < num_worker_threads; i++) {
    IDs[i] = i;
    pthread_create(&threads[i], NULL, worker_thread_func, &IDs[i]);
    set_thread_affinity(threads[i], cpus[i % num_cpus]);
  }

  for (i = 0; i < num_worker_threads; i++) {
    pthread_join(threads[i], NULL);
  }
  return 0;
}

