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
    return Row(
      children: [
        Padding(
          padding: const EdgeInsets.all(4.0),
          child: CircleAvatar(
            backgroundImage: (userAvatar != null)
                ? NetworkImage(userAvatar!)
                : AssetImage(Constants.defaultAvatarPng) as ImageProvider,
            radius: 25,
          ),
        ),
        Column(
          children: [
            Row(
              children: [
                Text(
                  userName,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Text(" "),
                Text(
                  activity,
                  style: TextStyle(
                    fontSize: 16,
                  ),
                )
              ],
            ),
            Text(
              postTime,
              style: TextStyle(
                color: Colors.grey,
              ),
            )
          ],
        ),
        IconButton(onPressed: () {}, icon: Icon(Icons.more_horiz))
      ],
    );
  }
}
