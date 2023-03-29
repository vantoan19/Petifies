// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/features/comment/controller/comment_count_controller.dart';
import 'package:mobile/src/features/comment/screens/create_comment_screen.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/posts/post.dart';

class CommentButton extends StatelessWidget {
  final double textSize;
  final double iconSize;
  final double spaceBetween;

  const CommentButton({
    Key? key,
    required this.textSize,
    required this.iconSize,
    required this.spaceBetween,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        CommentReplyButton(
          iconSize: iconSize,
        ),
        CommentCount(
          textStyle: TextStyle(
            fontSize: textSize,
            fontWeight: FontWeight.w300,
          ),
          leftPadding: spaceBetween,
        ),
      ],
    );
  }
}

class CommentReplyButton extends ConsumerWidget {
  final double iconSize;

  const CommentReplyButton({
    Key? key,
    required this.iconSize,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isPostTarget = ref.watch(isPostContextProvider);
    String targetID;

    if (isPostTarget) {
      targetID = ref.watch(postInfoProvider.select((info) => info.id));
    } else {
      targetID = ref.watch(commentInfoProvider.select((info) => info.id));
    }

    return IconButton(
      constraints: BoxConstraints(
        minWidth: 10,
        minHeight: 10,
      ),
      padding: EdgeInsets.all(0),
      onPressed: () {
        if (isPostTarget) {
          Navigator.pushNamed(
            context,
            '/post-details',
            arguments: Tuple2(ref.read(postInfoProvider), true),
          );
        }
      },
      icon: Icon(
        FontAwesomeIcons.comment,
        size: iconSize,
        color: Theme.of(context).colorScheme.secondary,
      ),
    );
  }
}

class CommentCount extends ConsumerWidget {
  final TextStyle textStyle;
  final double leftPadding;

  const CommentCount({
    Key? key,
    required this.textStyle,
    required this.leftPadding,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isPostTarget = ref.watch(isPostContextProvider);

    String targetID;
    int initialCount;
    if (isPostTarget) {
      targetID = ref.watch(postInfoProvider.select((info) => info.id));
      initialCount =
          ref.watch(postInfoProvider.select((info) => info.commentCount));
    } else {
      targetID = ref.watch(commentInfoProvider.select((info) => info.id));
      initialCount =
          ref.watch(commentInfoProvider.select((info) => info.subcommentCount));
    }

    ref.watch(commentCountControllerProvider(Tuple2(targetID, isPostTarget)));
    final commentCount =
        ref.watch(commentCountProvider(Tuple2(targetID, isPostTarget)));
    final hasChangedCommentCount = ref
        .watch(hasChangedCommentCountProvider(Tuple2(targetID, isPostTarget)));

    final count = (hasChangedCommentCount) ? commentCount : initialCount;

    return Padding(
      padding: EdgeInsets.fromLTRB(leftPadding, 0, 0, 0),
      child: Text(
        StringUtils.stringifyCount(count),
        style: textStyle,
      ),
    );
  }
}
