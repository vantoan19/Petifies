import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';

class UserAvatar extends StatelessWidget {
  final String? userAvatar;
  final EdgeInsetsGeometry? padding;

  const UserAvatar({Key? key, this.userAvatar = null, this.padding = null})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding:
          padding == null ? const EdgeInsets.fromLTRB(0, 0, 0, 0) : padding!,
      child: (userAvatar != null && userAvatar != "")
          ? CircleAvatar(
              backgroundImage: NetworkImage(userAvatar!),
              radius: 20,
              backgroundColor: Colors.transparent,
            )
          : CircleAvatar(
              backgroundImage: AssetImage(Constants.defaultAvatarPng),
              radius: 20,
              backgroundColor: Colors.transparent,
            ),
    );
  }
}
