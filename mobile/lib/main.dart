import 'package:flutter/material.dart';
import 'package:flutter_gen/gen_l10n/petifies_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/routers.dart';
import 'package:mobile/src/features/auth/screens/introduction_screens.dart';
import 'package:mobile/src/features/home/screens/home_screen.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/theme/themes.dart';

void main() {
  runApp(const ProviderScope(child: const Petifies()));
}

class Petifies extends ConsumerWidget {
  const Petifies({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    final isLoading = user.maybeWhen(
      data: (_) => user.isRefreshing,
      loading: () => true,
      orElse: () => false,
    );

    return MaterialApp(
      title: 'Petifies',
      localizationsDelegates: AppLocalizations.localizationsDelegates,
      supportedLocales: AppLocalizations.supportedLocales,
      locale: Locale('en', ''),
      theme: Themes.darkModeAppTheme,
      home: isLoading
          ? Center(
              child: CircularProgressIndicator(),
            )
          : userInfo == null
              ? IntroductionScreens()
              : IntroductionScreens(),
      // home: SignInScreen(),
      onGenerateRoute: onGenerateRoute,
    );
  }
}
