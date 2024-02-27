#ifndef DEPGRAPH_GO
#define DEPGRAPH_GO

#include <python3.11/Python.h>
#include <stdlib.h>
#include <string.h>

char **get_pkg_depgraph(const char *pkg_name);
void cleanup();

#endif
