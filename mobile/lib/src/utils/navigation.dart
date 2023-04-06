import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/comment/screens/create_comment_screen.dart';
import 'package:mobile/src/features/home/screens/home_screen.dart';
import 'package:mobile/src/features/media/screens/media_full_page_screen.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';

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

  static void toHomePage(context) {
    Navigator.of(context).pushNamed('/home-page');
  }

  static void toCreatePost(context) {
    Navigator.of(context, rootNavigator: true).pushNamed('/create-post');
  }

  static void goBack(context) {
    if (Navigator.of(context).canPop()) {
      Navigator.of(context).pop();
    }
  }

  static Widget showMediaFullPageBottomSheet({
    required WidgetRef ref,
    required String mediaUrl,
    required bool isMediaImage,
  }) {
    final isPostTarget = ref.read(isPostContextProvider);
    final postInfo = ref.read(postInfoProvider);
    return ProviderScope(
      overrides: [
        isPostContextProvider.overrideWithValue(isPostTarget),
        postInfoProvider.overrideWithValue(postInfo),
        if (!isPostTarget)
          commentInfoProvider.overrideWithValue(ref.read(commentInfoProvider)),
      ],
      child: MediaFullPageScreen(
        mediaUrl: mediaUrl,
        isMediaImage: isMediaImage,
      ),
    );
  }

  static Widget showCreateCommentBottomSheet({
    required WidgetRef ref,
  }) {
    final isPostTarget = ref.read(isPostContextProvider);
    final postInfo = ref.read(postInfoProvider);
    return ProviderScope(
      overrides: [
        isPostContextProvider.overrideWithValue(isPostTarget),
        postInfoProvider.overrideWithValue(postInfo),
        if (!isPostTarget)
          commentInfoProvider.overrideWithValue(ref.read(commentInfoProvider)),
      ],
      child: CreateCommentScreen(),
    );
  }
}
