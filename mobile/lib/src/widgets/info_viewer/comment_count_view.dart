// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/features/comment/controller/comment_count_controller.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/widgets/info_viewer/general_count_view.dart';

class CommentCountView extends ConsumerWidget {
  final double fontSize;
  const CommentCountView({
    required this.fontSize,
  });

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

    return GeneralCountView(
      count: count,
      singularLabel: "Comment",
      pluralLabel: "Comments",
      fontSize: fontSize,
      onTapCallback: () {},
    );
  }
}
