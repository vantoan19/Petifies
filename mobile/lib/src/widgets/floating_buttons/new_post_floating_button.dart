import 'package:flutter/material.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';

class NewPostFloatingButton extends StatelessWidget {
  const NewPostFloatingButton({super.key});

  @override
  Widget build(BuildContext context) {
    return FloatingActionButton(
      onPressed: () {
        NavigatorUtil.toCreatePost(context);
      },
      child: Icon(
        Icons.add,
        color: Theme.of(context).colorScheme.inversePrimary,
      ),
      backgroundColor: Themes.blueColor,
    );
  }
}
