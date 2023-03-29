import 'package:intl/intl.dart';

class StringUtils {
  static String stringifyCount(int count) {
    if (count < 10000) {
      return count.toString();
    } else if (count < 1000000) {
      String count_ = (count.toDouble() / 1000.0).toStringAsFixed(1);
      return '${count_}K';
    } else if (count < 1000000000) {
      String count_ = (count.toDouble() / 1000000.0).toStringAsFixed(1);
      return '${count_}M';
    } else {
      String count_ = (count.toDouble() / 1000000000.0).toStringAsFixed(1);
      return '${count_}B';
    }
  }

  static String stringifyTime(DateTime time) {
    final diff = DateTime.now().difference(time);
    if (diff.compareTo(Duration(seconds: 10)) < 0) {
      return "just now";
    } else if (diff.compareTo(Duration(minutes: 1)) < 0) {
      return diff.inSeconds.toString() + "s";
    } else if (diff.compareTo(Duration(hours: 1)) < 0) {
      return diff.inMinutes.toString() + "m";
    } else if (diff.compareTo(Duration(days: 1)) < 0) {
      return diff.inHours.toString() + "h";
    } else if (diff.compareTo(Duration(days: 30)) < 0) {
      return diff.inDays.toString() + "d";
    }
    final DateFormat formatter = DateFormat("dd-MM-yyyy");
    return formatter.format(time);
  }

  static String getActivity(String postActivity) {
    switch (postActivity) {
      case "post":
        return "shared a new post";
      default:
        return "shared a new post";
    }
  }
}
