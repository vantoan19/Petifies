// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/appbars/only_go_back_appbar.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';
import 'package:mobile/src/widgets/images/full_page_image.dart';
import 'package:mobile/src/widgets/textfield/reply_text_field.dart';
import 'package:mobile/src/widgets/videos/full_page_video_player.dart';

class MediaFullPageScreen extends StatelessWidget {
  final String mediaUrl;
  final bool isMediaImage;

  const MediaFullPageScreen({
    Key? key,
    required this.mediaUrl,
    required this.isMediaImage,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const OnlyGoBackAppbar(),
      body: Column(
        children: [
          Expanded(
            child: isMediaImage
                ? FullPageImageViewer(imageUrl: mediaUrl)
                : FullPageVideoPlayer(videoUrl: mediaUrl),
          ),
          const _ActionButtons(),
          const ReplyTextField(
            autoFocus: false,
          ),
        ],
      ),
      bottomNavigationBar: SizedBox(
        height: 16,
      ),
    );
  }
}

class _ActionButtons extends StatelessWidget {
  const _ActionButtons({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 20),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.heart),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 20,
          ),
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.comment),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 20,
          ),
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.retweet),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 20,
          ),
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.paperPlane),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 20,
          ),
        ],
      ),
    );
  }
}
