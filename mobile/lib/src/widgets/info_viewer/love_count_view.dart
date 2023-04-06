// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/features/love/controllers/love_count_controller.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/widgets/info_viewer/general_count_view.dart';

class LoveCountView extends ConsumerWidget {
  final double fontSize;
  const LoveCountView({
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
          ref.watch(postInfoProvider.select((info) => info.loveCount));
    } else {
      targetID = ref.watch(commentInfoProvider.select((info) => info.id));
      initialCount =
          ref.watch(commentInfoProvider.select((info) => info.loveCount));
    }

    ref.watch(loveCountControllerProvider(Tuple2(targetID, isPostTarget)));
    final loveCount =
        ref.watch(loveCountProvider(Tuple2(targetID, isPostTarget)));
    final hasChangedLoveCount =
        ref.watch(hasChangedLoveCountProvider(Tuple2(targetID, isPostTarget)));

    final count = (hasChangedLoveCount) ? loveCount : initialCount;

    return GeneralCountView(
      count: count,
      singularLabel: "Love",
      pluralLabel: "Loves",
      fontSize: fontSize,
      onTapCallback: () {},
    );
  }
}
