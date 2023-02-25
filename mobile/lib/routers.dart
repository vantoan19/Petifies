import 'package:flutter/material.dart';
import 'package:mobile/src/features/auth/screens/introduction_screens.dart';
import 'package:mobile/src/features/auth/screens/signin_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_form_screen.dart';
import 'package:mobile/src/features/auth/screens/splash_screen.dart';
import 'package:mobile/src/features/feed/screens/home_screen.dart';

Route onGenerateRoute(RouteSettings settings) {
  switch (settings.name) {
    case "/splash":
      return MaterialPageRoute(builder: (context) => const SplashScreen());
    case "/introduction":
      return MaterialPageRoute(
          builder: (context) => const IntroductionScreens());
    case "/signin":
      return MaterialPageRoute(builder: (context) => const SignInScreen());
    case "/signup":
      return MaterialPageRoute(builder: (context) => const SignInScreen());
    case "/signup/form":
      return MaterialPageRoute(builder: (context) => const SignUpFormScreen());
    case "/home-page":
      return MaterialPageRoute(builder: (context) => const HomeScreeen());
    default:
      return MaterialPageRoute(builder: (context) => const SignInScreen());
  }
}
