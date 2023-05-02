import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/petifies_proposal.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/status/status.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';

class Review extends ConsumerWidget {
  const Review({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final author = ref.read(reviewProvider.select((value) => value.author));
    final review = ref.read(reviewProvider.select((value) => value.review));
    final createdAt =
        ref.read(reviewProvider.select((value) => value.createdAt));

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            UserAvatar(
              userAvatar: author.userAvatar,
              padding: EdgeInsets.only(right: 16),
            ),
            Column(
              mainAxisSize: MainAxisSize.max,
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  author.firstName + " " + author.lastName,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Theme.of(context).colorScheme.secondary,
                  ),
                ),
                Text(
                  StringUtils.stringifyTime(createdAt),
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w300,
                    color: Theme.of(context)
                        .colorScheme
                        .secondary
                        .withOpacity(0.8),
                  ),
                ),
              ],
            )
          ],
        ),
        Padding(
          padding: const EdgeInsets.only(top: 14.0),
          child: Text(
            review,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w400,
              color: Theme.of(context).colorScheme.secondary,
              height: 1.35,
            ),
          ),
        ),
      ],
    );
  }
}
