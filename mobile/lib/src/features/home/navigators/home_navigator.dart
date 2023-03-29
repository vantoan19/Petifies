// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/post/screens/post_detail_screen.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/posts/post.dart';

class HomeNavigatorRoutes {
  static const String root = '/';
  static const String postDetail = '/post-details';
}

class HomeNavigator extends StatelessWidget {
  final GlobalKey<NavigatorState> navigatorKey;

  const HomeNavigator({
    Key? key,
    required this.navigatorKey,
  }) : super(key: key);

  // void _pushPostDetails(BuildContext context, PostModel postData) {
  //   Navigator.pushNamed(
  //     context,
  //     '/post-details',
  //     arguments: PostDetailsScreenArguments(postData: postData),
  //   );
  // }

  @override
  Widget build(BuildContext context) {
    return Navigator(
      key: navigatorKey,
      initialRoute: HomeNavigatorRoutes.root,
      onGenerateRoute: onGenerateRoute,
    );
  }

  Route onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case HomeNavigatorRoutes.root:
        return MaterialPageRoute(builder: (context) => FeedScreen());
      case HomeNavigatorRoutes.postDetail:
        final args = settings.arguments as Tuple2<PostModel, bool>;
        return MaterialPageRoute(
          builder: (context) => ProviderScope(
            overrides: [
              postInfoProvider.overrideWith((ref) => args.first),
              isPostContextProvider.overrideWithValue(true),
            ],
            child: PostDetailScreen(
              autoFocus: args.second,
            ),
          ),
        );
      default:
        return MaterialPageRoute(builder: (context) => FeedScreen());
    }
  }
}
