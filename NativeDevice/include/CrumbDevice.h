#pragma once
#include <stdint.h>
// Jobs own their workers. free cancels and releases the caller's reference;
// workers retain their own reference until all callbacks have stopped.
void *crumb_device_start(const char *json);
char *crumb_device_poll(void *job);
void crumb_device_cancel(void *job);
void crumb_device_free(void *job);
void crumb_device_string_free(char *value);
