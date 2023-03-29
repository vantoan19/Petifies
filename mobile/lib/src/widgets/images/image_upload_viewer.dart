// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:mobile/src/theme/themes.dart';

class ImageUploadViewer extends StatelessWidget {
  final File image;
  final String uploaderID;
  final VoidCallback removeAction;

  const ImageUploadViewer({
    Key? key,
    required this.image,
    required this.uploaderID,
    required this.removeAction,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 6, 16, 6),
      child: ClipRRect(
        child: SizedBox(
          height: MediaQuery.of(context).size.width * 0.5,
          child: Stack(
            alignment: Alignment.topRight,
            children: [
              Image.file(
                image,
                fit: BoxFit.contain,
              ),
              GestureDetector(
                onTap: () => removeAction(),
                child: Padding(
                  padding: const EdgeInsets.all(12),
                  child: Container(
                    width: 23,
                    height: 23,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: Themes.whiteColor,
                    ),
                    child: Icon(
                      Icons.cancel,
                      color: Themes.blackColor,
                      size: 22,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
        borderRadius: const BorderRadius.all(Radius.circular(12)),
      ),
    );
  }
}
