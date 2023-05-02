// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/navigators/explore_navigator.dart';
import 'package:mobile/src/features/petifies/navigators/my_petifies_navigator.dart';
import 'package:mobile/src/features/petifies/screens/my_petifies_screen.dart';
import 'package:mobile/src/widgets/appbars/go_back_title_appbar.dart';

class PetifiesDashboard extends StatelessWidget {
  const PetifiesDashboard({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const GoBackTitleAppbar(title: "Petifies dashboard"),
      body: Padding(
        padding: const EdgeInsets.fromLTRB(
          Constants.petifiesExpoloreHorizontalPadding,
          0,
          Constants.petifiesExpoloreHorizontalPadding,
          0,
        ),
        child: SingleChildScrollView(
          child: Column(
            children: [
              const SizedBox(
                height: 40,
              ),
              _Selection(
                action: () {
                  Navigator.of(context)
                      .pushNamed(PetifiesDashboardRoutes.myPetifies);
                },
                icon: Icons.pets,
                label: "My Petifies",
              ),
              const SizedBox(
                height: 16,
              ),
              _Selection(
                action: () {
                  Navigator.of(context)
                      .pushNamed(PetifiesDashboardRoutes.myProposals);
                },
                icon: FontAwesomeIcons.fileSignature,
                label: "My Proposals",
              ),
              const SizedBox(
                height: 16,
              ),
              _Selection(
                action: () {
                  Navigator.of(context)
                      .pushNamed(PetifiesDashboardRoutes.myReviews);
                },
                icon: Icons.rate_review_outlined,
                label: "My Reviews",
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _Selection extends StatelessWidget {
  final VoidCallback action;
  final IconData icon;
  final String label;
  const _Selection({
    Key? key,
    required this.action,
    required this.icon,
    required this.label,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return OutlinedButton(
      onPressed: action,
      child: Row(children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20),
          child: Icon(
            icon,
            size: 36,
            color: Theme.of(context).colorScheme.inversePrimary,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            fontSize: 16,
            color: Theme.of(context).colorScheme.inversePrimary,
          ),
        ),
      ]),
      style: OutlinedButton.styleFrom(
        minimumSize: Size.fromHeight(80),
        padding: EdgeInsets.zero,
        side: BorderSide(
          color: Theme.of(context).colorScheme.secondary.withOpacity(0.3),
          width: 2,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    );
  }
}
