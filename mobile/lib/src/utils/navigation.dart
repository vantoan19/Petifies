import 'package:flutter/material.dart';

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

  static void goBack(context) {
    if (Navigator.of(context).canPop()) {
      Navigator.of(context).pop();
    }
  }
}
