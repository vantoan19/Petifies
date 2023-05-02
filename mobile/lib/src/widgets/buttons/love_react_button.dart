// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/features/love/controllers/love_count_controller.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';

class LoveReactButton extends StatelessWidget {
  final double iconSize;
  final double textSize;
  final double spaceBetween;

  const LoveReactButton({
    Key? key,
    required this.iconSize,
    required this.textSize,
    required this.spaceBetween,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        LoveReactIconButton(
          iconSize: iconSize,
        ),
        LoveCount(
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

class LoveReactIconButton extends ConsumerWidget {
  final double iconSize;

  LoveReactIconButton({
    required this.iconSize,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isPostTarget = ref.watch(isPostContextProvider);

    String targetID;
    bool initialHasReacted;

    if (isPostTarget) {
      targetID = ref.watch(postInfoProvider.select((info) => info.id));
      initialHasReacted =
          ref.watch(postInfoProvider.select((info) => info.hasReacted));
    } else {
      targetID = ref.watch(commentInfoProvider.select((info) => info.id));
      initialHasReacted =
          ref.watch(commentInfoProvider.select((info) => info.hasReacted));
    }

    final hasReacted =
        ref.watch(hasReactedProvider(Tuple2(targetID, isPostTarget)));

    final solidHeartIcon = Icon(
      FontAwesomeIcons.solidHeart,
      color: Themes.blueColor,
      size: iconSize,
    );

    final emptyHeartIcon = Icon(
      FontAwesomeIcons.heart,
      color: Theme.of(context).colorScheme.secondary,
      size: iconSize,
    );

    return NoPaddingIconButton(
      onPressed: () {
        ref
            .read(loveCountControllerProvider(Tuple2(targetID, isPostTarget))
                .notifier)
            .toggleLoveReact();
      },
      icon: (hasReacted == null)
          ? (initialHasReacted ? solidHeartIcon : emptyHeartIcon)
          : (hasReacted ? solidHeartIcon : emptyHeartIcon),
    );
  }
}

class LoveCount extends ConsumerWidget {
  final TextStyle textStyle;
  final double leftPadding;

  LoveCount({
    required this.textStyle,
    required this.leftPadding,
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

    return Padding(
      padding: EdgeInsets.fromLTRB(leftPadding, 0, 0, 0),
      child: Text(
        StringUtils.stringifyCount(count),
        style: textStyle,
      ),
    );
  }
}
