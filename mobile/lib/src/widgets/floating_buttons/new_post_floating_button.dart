// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';
import 'package:mobile/src/features/post/screens/create_post_screen.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';

class NewPostFloatingButton extends StatelessWidget {
  final String heroTag;
  const NewPostFloatingButton({
    Key? key,
    required this.heroTag,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return FloatingActionButton(
      heroTag: heroTag,
      onPressed: () {
        showModalBottomSheet(
            context: context,
            isScrollControlled: true,
            useSafeArea: true,
            barrierColor: Theme.of(context).scaffoldBackgroundColor,
            builder: (context) {
              return const CreatePostScreen();
            });
      },
      child: Icon(
        Icons.add,
        color: Theme.of(context).colorScheme.inversePrimary,
      ),
      backgroundColor: Themes.blueColor,
    );
  }
}
