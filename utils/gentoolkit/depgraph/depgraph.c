#include "./depgraph.h"

const char *stmt_pre_pkg_name = "from gentoolkit.dependencies import "
                                "Dependencies; result = [dep.__str__() "
                                "for dep in Dependencies('";
const char *stmt_post_pkg_name = "').get_all_depends()]";
PyObject *globals, *locals, *result;

void cleanup() {
  if (result != NULL)
    Py_DECREF(result);
  if (result != NULL)
    Py_DECREF(locals);
  if (result != NULL)
    Py_DECREF(globals);
  Py_Finalize();
}

char **get_pkg_depgraph(const char *pkg_name) {
  Py_Initialize();

  const int stmt_len =
      strlen(stmt_pre_pkg_name) + strlen(stmt_post_pkg_name) + strlen(pkg_name);
  char stmt[stmt_len];

  strcat(stmt, stmt_pre_pkg_name);
  strcat(stmt, pkg_name);
  strcat(stmt, stmt_post_pkg_name);

  globals = PyDict_New();
  locals = PyDict_New();
  result = PyRun_String(stmt, Py_file_input, globals, locals);

  if (result == NULL)
    return NULL;

  PyObject *value = PyDict_GetItem(locals, PyUnicode_FromString("result"));

  if (value == NULL)
    return NULL;

  int size = PyList_Size(value);
  char **deps = (char **)malloc((size + 1) * sizeof(char *));

  for (int i = 0; i < size; i++) {
    PyObject *item = PyList_GetItem(value, i);

    if (PyUnicode_Check(item)) {
      char *dep = (char *)PyUnicode_AsUTF8(item);
      deps[i] = dep;
    }
  }
  // to determine where to stop reading in go
  deps[size] = NULL;

  return deps;
}
