#include "./depgraph.h"

const char *get_deps_stmt_format =
    "from gentoolkit.dependencies import "
    "Dependencies; result = [dep.__str__() "
    "for dep in Dependencies('%s').get_all_depends()]";
PyObject *globals, *locals;

char *get_stmt_with_pkg_name(const char *pkg_name) {
  char *res = malloc(strlen(get_deps_stmt_format) + strlen(pkg_name));
  sprintf(res, get_deps_stmt_format, pkg_name);
  return res;
}

char **get_deps_from_python_result(PyObject *result) {
  int size = PyList_Size(result);
  char **deps = (char **)malloc((size + 1) * sizeof(char *));

  for (int i = 0; i < size; i++) {
    PyObject *item = PyList_GetItem(result, i);

    if (PyUnicode_Check(item)) {
      char *dep = (char *)PyUnicode_AsUTF8(item);
      deps[i] = dep;
    }
  }
  // to determine where to stop reading in go
  deps[size] = NULL;

  return deps;
}

char **get_pkg_depgraph(const char *pkg_name) {
  Py_Initialize();

  globals = PyDict_New();
  locals = PyDict_New();
  PyRun_String(get_stmt_with_pkg_name(pkg_name), Py_file_input, globals,
               locals);

  PyObject *result = PyDict_GetItem(locals, PyUnicode_FromString("result"));

  if (result == NULL)
    return NULL;

  return get_deps_from_python_result(result);
}

void cleanup() {
  if (locals != NULL)
    Py_DECREF(locals);
  if (globals != NULL)
    Py_DECREF(globals);
  Py_Finalize();
}
