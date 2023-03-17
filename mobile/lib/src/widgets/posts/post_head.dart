import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';

class PostHead extends StatelessWidget {
  final String? userAvatar;
  final String userName;
  final String activity;
  final String postTime;

  const PostHead({
    super.key,
    this.userAvatar = null,
    required this.userName,
    required this.activity,
    required this.postTime,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 8, 0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          // Avatar
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 0, 0, 0),
            child: (userAvatar != null)
                ? CircleAvatar(
                    backgroundImage: NetworkImage(userAvatar!),
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
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(0, 4, 0, 4),
                child: Row(
                  children: [
                    // Name
                    Text(
                      userName,
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const Text(" "),
                    // Activity
                    Text(
                      activity,
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
                postTime,
                style: TextStyle(
                  color: Colors.grey,
                  fontWeight: FontWeight.w300,
                ),
              )
            ],
          ),
          // More button
          IconButton(onPressed: () {}, icon: Icon(Icons.more_horiz))
        ],
      ),
    );
  }
}
