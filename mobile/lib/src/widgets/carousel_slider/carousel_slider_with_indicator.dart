// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:carousel_slider/carousel_controller.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:page_view_dot_indicator/page_view_dot_indicator.dart';

class ImageCarouselSliderWithIndicators extends StatefulWidget {
  final List<NetworkImageModel> images;
  final double height;
  final double width;
  final bool roundCorner;
  final EdgeInsetsGeometry? padding;

  const ImageCarouselSliderWithIndicators({
    Key? key,
    required this.images,
    required this.height,
    required this.width,
    required this.roundCorner,
    this.padding = null,
  }) : super(key: key);

  @override
  State<ImageCarouselSliderWithIndicators> createState() =>
      _ImageCarouselSliderWithIndicatorsState();
}

class _ImageCarouselSliderWithIndicatorsState
    extends State<ImageCarouselSliderWithIndicators> {
  int _currentIndex = 0;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: widget.width,
      height: widget.height,
      padding: widget.padding,
      child: ClipRRect(
        borderRadius: BorderRadius.only(
          topLeft: widget.roundCorner ? Radius.circular(16) : Radius.zero,
          topRight: widget.roundCorner ? Radius.circular(16) : Radius.zero,
          bottomLeft: widget.roundCorner ? Radius.circular(16) : Radius.zero,
          bottomRight: widget.roundCorner ? Radius.circular(16) : Radius.zero,
        ),
        child: Stack(
          fit: StackFit.expand,
          children: [
            CarouselSlider(
              items: widget.images
                  .map((e) => Image.network(
                        e.uri,
                        fit: BoxFit.cover,
                        width: widget.width,
                      ))
                  .toList(),
              options: CarouselOptions(
                  aspectRatio:
                      widget.width.toDouble() / widget.height.toDouble(),
                  autoPlay: false,
                  enlargeCenterPage: true,
                  viewportFraction: 1,
                  onPageChanged: (index, reason) {
                    setState(() {
                      _currentIndex = index;
                    });
                  }),
            ),
            Align(
              alignment: Alignment.bottomCenter,
              child: Container(
                width: 150,
                padding: const EdgeInsets.symmetric(
                  vertical: 16,
                  horizontal: 24,
                ),
                child: PageViewDotIndicator(
                  currentItem: _currentIndex,
                  count: widget.images.length,
                  unselectedColor: Theme.of(context).colorScheme.secondary,
                  selectedColor: Themes.blueColor,
                  size: const Size(8, 8),
                  unselectedSize: const Size(8, 8),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
