// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter/src/widgets/placeholder.dart';

class FullPageImageViewer extends StatelessWidget {
  final String imageUrl;
  const FullPageImageViewer({
    Key? key,
    required this.imageUrl,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Image.network(
      imageUrl,
      loadingBuilder: (context, child, loadingProgress) =>
          (loadingProgress == null) ? child : CircularProgressIndicator(),
      fit: BoxFit.contain,
    );
  }
}
