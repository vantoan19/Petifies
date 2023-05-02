import 'package:flutter/material.dart';
import 'package:mobile/src/features/auth/screens/introduction_screens.dart';
import 'package:mobile/src/features/auth/screens/signin_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_form_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_screen.dart';
import 'package:mobile/src/features/auth/screens/splash_screen.dart';
import 'package:mobile/src/features/comment/screens/create_comment_screen.dart';
import 'package:mobile/src/features/home/screens/home_screen.dart';
import 'package:mobile/src/features/media/screens/media_full_page_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_home_screen.dart';
import 'package:mobile/src/features/post/screens/create_post_screen.dart';

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
      return MaterialPageRoute(builder: (context) => const SignUpScreen());
    case "/signup/form":
      return MaterialPageRoute(builder: (context) => const SignUpFormScreen());
    case "/home-page":
      return MaterialPageRoute(builder: (context) => const HomeScreeen());
    case "/create-post":
      return MaterialPageRoute(builder: (context) => const CreatePostScreen());
    case "/petifies-home":
      return MaterialPageRoute(
          builder: (context) => const PetifiesHomeScreeen());
    case "/create-petifies":
      return MaterialPageRoute(
          builder: (context) => const CreatePetifiesScreen());
    default:
      return MaterialPageRoute(builder: (context) => const SignInScreen());
  }
}
