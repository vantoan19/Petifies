import 'package:flutter/material.dart';
import 'package:mobile/src/features/auth/screens/introduction_screens.dart';
import 'package:mobile/src/features/auth/screens/signin_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_form_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_screen.dart';
import 'package:mobile/src/features/auth/screens/splash_screen.dart';
import 'package:mobile/src/theme/themes.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Petifies',
      theme: Themes.lightModeAppTheme,
      initialRoute: "/introduction",
      onGenerateRoute: (settings) {
        switch (settings.name) {
          case "/splash":
            return MaterialPageRoute(
                builder: (context) => const SplashScreen());
          case "/introduction":
            return MaterialPageRoute(
                builder: (context) => const IntroductionScreens());
          case "/signin":
            return MaterialPageRoute(
                builder: (context) => const SignInScreen());
          case "/signup":
            return MaterialPageRoute(
                builder: (context) => const SignUpScreen());
          case "/signup/form":
            return MaterialPageRoute(
                builder: (context) => const SignUpFormScreen());
          default:
            return MaterialPageRoute(
                builder: (context) => const SignInScreen());
        }
      },
    );
  }
}
