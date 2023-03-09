#ifndef DEPGRAPH_GO
#define DEPGRAPH_GO

#include <python3.10/Python.h>
#include <string.h>
#include <stdlib.h>

char **get_pkg_depgraph(const char *pkg_name);
void cleanup();

#endif
