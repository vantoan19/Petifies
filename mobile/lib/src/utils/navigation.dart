import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/comment/screens/create_comment_screen.dart';
import 'package:mobile/src/features/home/screens/home_screen.dart';
import 'package:mobile/src/features/media/screens/media_full_page_screen.dart';
import 'package:mobile/src/models/basic_user_info.dart';

class NavigatorUtil {
  static void toSignIn(context) {
    Navigator.of(context).pushNamed('/signin');
  }

  static void toSignUp(context) {
    Navigator.of(context).pushNamed('/signup');
  }

  static void toSignUpForm(context) {
    Navigator.of(context).pushNamed('/signup/form');
  }

  static void toIntroductionPages(context) {
    Navigator.of(context).pushNamed('/introduction');
  }

  static void toHomePage(context) {
    Navigator.of(context).pushNamed('/home-page');
  }

  static void toCreatePost(context) {
    Navigator.of(context, rootNavigator: true).pushNamed('/create-post');
  }

  static void goBack(context) {
    if (Navigator.of(context).canPop()) {
      Navigator.of(context).pop();
    }
  }
}
