#include "./depgraph.h"

PyObject *globals, *locals, *result;

void cleanup() {
  Py_DECREF(result);
  Py_DECREF(locals);
  Py_DECREF(globals);
  Py_Finalize();
}

const char *get_pkg_depgraph(const char *pkg_name) {
  Py_Initialize();

  PyRun_SimpleString("from gentoolkit.dependencies import Dependencies");
  const char *stmt_pre_pkg_name =
      "from gentoolkit.dependencies import Dependencies; result = \" "
      "\".join(Dependencies('";
  const char *stmt_post_pkg_name = "').get_all_depends())";

  const int stmt_len =
      strlen(stmt_pre_pkg_name) + strlen(stmt_post_pkg_name) + strlen(pkg_name);
  char stmt[stmt_len];

  strcat(stmt, stmt_pre_pkg_name);
  strcat(stmt, pkg_name);
  strcat(stmt, stmt_post_pkg_name);

  globals = PyDict_New();
  locals = PyDict_New();
  result = PyRun_String(stmt, Py_file_input, globals, locals);

  if (result == NULL) {
    return NULL;
  }

  PyObject *value = PyDict_GetItemString(locals, "result");

  if (value == NULL) {
    return NULL;
  }

  return PyUnicode_AsUTF8(value);
}
