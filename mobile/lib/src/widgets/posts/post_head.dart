// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/posts/uploading_post.dart';

class PostHead extends ConsumerWidget {
  final bool isUploadingPost;

  const PostHead({
    required this.isUploadingPost,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    BasicUserInfoModel owner;
    DateTime createdAt;
    String activity;

    if (isUploadingPost) {
      owner = ref.watch(uploadingPostInfoProvider.select((info) => info.owner));
      createdAt =
          ref.watch(uploadingPostInfoProvider.select((info) => info.createdAt));
      activity = ref
          .watch(uploadingPostInfoProvider.select((info) => info.postActivity));
    } else {
      owner = ref.watch(postInfoProvider.select((info) => info.owner));
      createdAt = ref.watch(postInfoProvider.select((info) => info.createdAt));
      activity =
          ref.watch(postInfoProvider.select((info) => info.postActivity));
    }

    return Padding(
      padding: const EdgeInsets.fromLTRB(
        Constants.horizontalScreenPadding,
        0,
        Constants.horizontalScreenPadding,
        0,
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          // Avatar
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 0, 0),
            child: (owner.userAvatar != null)
                ? CircleAvatar(
                    backgroundImage: NetworkImage(owner.userAvatar!),
                    radius: 25,
                    backgroundColor: Colors.transparent,
                  )
                : CircleAvatar(
                    backgroundImage: AssetImage(Constants.defaultAvatarPng),
                    radius: 25,
                    backgroundColor: Colors.transparent,
                  ),
          ),
          // Name & activity
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 4, 0, 4),
                    child: Row(
                      children: [
                        // Name
                        Text(
                          owner.firstName + " " + owner.lastName,
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const Text(" "),
                        // Activity
                        Text(
                          StringUtils.getActivity(activity),
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w300,
                          ),
                        )
                      ],
                    ),
                  ),
                  // Time
                  Text(
                    StringUtils.stringifyTime(createdAt),
                    style: TextStyle(
                      color: Colors.grey,
                      fontWeight: FontWeight.w300,
                    ),
                  )
                ],
              ),
            ),
          ),
          // More button
          IconButton(
            onPressed: () {},
            icon: Icon(Icons.more_horiz),
            padding: EdgeInsets.zero,
            constraints: BoxConstraints(minHeight: 30, minWidth: 40),
            color: Theme.of(context).colorScheme.secondary,
          )
        ],
      ),
    );
  }
}
