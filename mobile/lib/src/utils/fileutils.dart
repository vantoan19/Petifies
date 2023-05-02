import 'dart:io';
import 'package:path/path.dart';

class FileUtils {
  static String getFilename(File file) {
    return basename(file.path);
  }

  static String getFileExtension(File file) {
    return extension(file.path);
  }

  static int getFileSize(File file) {
    return file.lengthSync();
  }
}
