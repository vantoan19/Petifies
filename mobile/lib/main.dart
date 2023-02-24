import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/petifies_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/auth/screens/introduction_screens.dart';
import 'package:mobile/src/features/auth/screens/signin_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_form_screen.dart';
import 'package:mobile/src/features/auth/screens/signup_screen.dart';
import 'package:mobile/src/features/auth/screens/splash_screen.dart';
import 'package:mobile/src/theme/themes.dart';

void main() {
  runApp(const ProviderScope(child: const Petifies()));
}

class Petifies extends StatelessWidget {
  const Petifies({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Petifies',
      localizationsDelegates: AppLocalizations.localizationsDelegates,
      supportedLocales: AppLocalizations.supportedLocales,
      locale: Locale('vi', ''),
      theme: Themes.darkModeAppTheme,
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
