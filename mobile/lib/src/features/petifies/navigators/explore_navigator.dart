// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/petifies/screens/my_petifies_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_details_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_explore_screen.dart';
import 'package:mobile/src/providers/petifies_providers.dart';

class PetifiesExploreNavigatorRoutes {
  static const String root = '/';
  static const String petifiesDetails = '/petifies-details';
}

class PetifiesExploreNavigator extends ConsumerWidget {
  final GlobalKey<NavigatorState> navigatorKey;

  const PetifiesExploreNavigator({
    Key? key,
    required this.navigatorKey,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Navigator(
      key: navigatorKey,
      initialRoute: PetifiesExploreNavigatorRoutes.root,
      onGenerateRoute: (settings) => onGenerateRoute(settings, ref),
    );
  }

  Route onGenerateRoute(RouteSettings settings, WidgetRef ref) {
    switch (settings.name) {
      case PetifiesExploreNavigatorRoutes.root:
        return MaterialPageRoute(
            builder: (context) => const PetifiesExploreScreen());
      case PetifiesExploreNavigatorRoutes.petifiesDetails:
        final args = settings.arguments as PetifiesDetailsScreenArguments;
        return MaterialPageRoute(
          builder: (context) => ProviderScope(
            overrides: [
              petifiesInfoProvider.overrideWithValue(args.petifiesData),
            ],
            child: const PetifiesDetailsScreen(),
          ),
        );
      default:
        return MaterialPageRoute(builder: (context) => FeedScreen());
    }
  }
}
