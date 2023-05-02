// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/theme/themes.dart';

class CustomStep extends StatelessWidget {
  final Widget leftWidget;
  final Widget rightColumnWidget;
  final bool showStepLine;
  final double spaceBetween;
  final EdgeInsetsGeometry? padding;

  const CustomStep({
    Key? key,
    required this.leftWidget,
    required this.rightColumnWidget,
    required this.showStepLine,
    required this.spaceBetween,
    this.padding,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding ?? EdgeInsets.zero,
      child: IntrinsicHeight(
        child: Row(
          children: [
            Padding(
              padding: EdgeInsets.only(right: spaceBetween),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.start,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  leftWidget,
                  if (showStepLine)
                    Expanded(
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.all(Radius.circular(100)),
                          color: Theme.of(context).colorScheme.secondary,
                        ),
                        margin: EdgeInsets.symmetric(vertical: 8),
                        width: 1.5,
                      ),
                    )
                ],
              ),
            ),
            Expanded(
              child: rightColumnWidget,
            )
          ],
        ),
      ),
    );
  }
}
