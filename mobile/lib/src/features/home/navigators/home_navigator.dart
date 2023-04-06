// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/features/comment/screens/comment_details_screen.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/post/screens/post_detail_screen.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/posts/post.dart';

class HomeNavigatorRoutes {
  static const String root = '/';
  static const String postDetails = '/post-details';
  static const String commentDetails = '/comment-details';
}

class HomeNavigator extends ConsumerWidget {
  final GlobalKey<NavigatorState> navigatorKey;

  const HomeNavigator({
    Key? key,
    required this.navigatorKey,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Navigator(
      key: navigatorKey,
      initialRoute: HomeNavigatorRoutes.root,
      onGenerateRoute: (settings) => onGenerateRoute(settings, ref),
    );
  }

  Route onGenerateRoute(RouteSettings settings, WidgetRef ref) {
    switch (settings.name) {
      case HomeNavigatorRoutes.root:
        return MaterialPageRoute(builder: (context) => FeedScreen());
      case HomeNavigatorRoutes.postDetails:
        final args = settings.arguments as PostDetailsScreenArguments;
        return MaterialPageRoute(
          builder: (context) => ProviderScope(
            overrides: [
              postInfoProvider.overrideWithValue(args.postData),
              isPostContextProvider.overrideWithValue(true),
            ],
            child: PostDetailScreen(
              autoFocus: args.autoFocus,
            ),
          ),
        );
      case HomeNavigatorRoutes.commentDetails:
        final args = settings.arguments as CommentDetailsScreenArguments;
        return MaterialPageRoute(
          builder: (context) => ProviderScope(
            overrides: [
              postInfoProvider.overrideWithValue(args.postData),
              isPostContextProvider.overrideWithValue(false),
              commentInfoProvider.overrideWithValue(args.commentData),
              ancestorCommentsProvider.overrideWithValue(args.ancestorComments),
            ],
            child: CommentDetailScreen(
              autoFocus: args.autoFocus,
            ),
          ),
        );
      default:
        return MaterialPageRoute(builder: (context) => FeedScreen());
    }
  }
}
