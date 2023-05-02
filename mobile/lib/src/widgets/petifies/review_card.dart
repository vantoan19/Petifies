import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';
import 'package:readmore/readmore.dart';

class ReviewCard extends ConsumerWidget {
  const ReviewCard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final author = ref.watch(reviewProvider.select((value) => value.author));
    final review = ref.watch(reviewProvider.select((value) => value.review));
    return Card(
      color: Theme.of(context).colorScheme.tertiary,
      child: Container(
        width: 250,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(16.0, 16, 16, 0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ListTile(
                contentPadding: EdgeInsets.zero,
                leading: UserAvatar(
                  userAvatar: author.userAvatar,
                ),
                title: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      author.firstName + " " + author.lastName,
                      style: TextStyle(
                        fontSize: 16,
                      ),
                    ),
                  ],
                ),
                subtitle: Padding(
                  padding: const EdgeInsets.only(top: 8),
                  child: ReadMoreText(
                    review,
                    trimMode: TrimMode.Line,
                    trimLines: 3,
                    trimCollapsedText: '',
                    trimExpandedText: '',
                  ),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  TextButton(
                    child: const Text(
                      'Show more',
                      style: TextStyle(
                        fontSize: 12,
                      ),
                    ),
                    onPressed: () {},
                    style: ButtonStyle(
                      padding: MaterialStatePropertyAll(EdgeInsets.zero),
                    ),
                  ),
                ],
              )
            ],
          ),
        ),
      ),
    );
  }
}
