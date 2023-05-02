// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/utils/stringutils.dart';

class GeneralCountView extends StatelessWidget {
  final int count;
  final String singularLabel;
  final String pluralLabel;
  final double fontSize;
  final VoidCallback onTapCallback;

  const GeneralCountView({
    Key? key,
    required this.count,
    required this.singularLabel,
    required this.pluralLabel,
    required this.fontSize,
    required this.onTapCallback,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTapCallback,
      child: Row(
        children: [
          Text(
            StringUtils.stringifyCount(count),
            style: TextStyle(
              fontSize: fontSize,
              fontWeight: FontWeight.w900,
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 0, 0, 0),
            child: Text(
              count < 2 ? singularLabel : pluralLabel,
              style: TextStyle(
                fontSize: fontSize,
                color: Theme.of(context).colorScheme.secondary,
                fontWeight: FontWeight.w300,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
