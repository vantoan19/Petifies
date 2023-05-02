// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/petifies/screens/my_petifies_screen.dart';
import 'package:mobile/src/features/petifies/screens/my_proposals_screen.dart';
import 'package:mobile/src/features/petifies/screens/my_reviews_screen.dart';
import 'package:mobile/src/features/petifies/screens/petifies_dashboard.dart';
import 'package:mobile/src/features/petifies/screens/petifies_details_screen.dart';
import 'package:mobile/src/providers/petifies_providers.dart';

class PetifiesDashboardRoutes {
  static const String root = '/';
  static const String petifiesDetails = '/petifies-details';
  static const String myPetifies = '/my-petifies';
  static const String myProposals = '/my-proposals';
  static const String myReviews = '/my-reviews';
}

class PetifiesDashboardNavigator extends ConsumerWidget {
  final GlobalKey<NavigatorState> navigatorKey;

  const PetifiesDashboardNavigator({
    Key? key,
    required this.navigatorKey,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Navigator(
      key: navigatorKey,
      initialRoute: PetifiesDashboardRoutes.root,
      onGenerateRoute: (settings) => onGenerateRoute(settings, ref),
    );
  }

  Route onGenerateRoute(RouteSettings settings, WidgetRef ref) {
    switch (settings.name) {
      case PetifiesDashboardRoutes.root:
        return MaterialPageRoute(
            builder: (context) => const PetifiesDashboard());
      case PetifiesDashboardRoutes.petifiesDetails:
        final args = settings.arguments as PetifiesDetailsScreenArguments;
        return MaterialPageRoute(
          builder: (context) => ProviderScope(
            overrides: [
              petifiesInfoProvider.overrideWithValue(args.petifiesData),
            ],
            child: const PetifiesDetailsScreen(),
          ),
        );
      case PetifiesDashboardRoutes.myPetifies:
        return MaterialPageRoute(
            builder: (context) => const MyPetifiesScreen());
      case PetifiesDashboardRoutes.myProposals:
        return MaterialPageRoute(
            builder: (context) => const MyProposalsScreen());
      case PetifiesDashboardRoutes.myReviews:
        return MaterialPageRoute(builder: (context) => const MyReviewScreen());
      default:
        return MaterialPageRoute(builder: (context) => FeedScreen());
    }
  }
}
