DX4WIN.CTY #22.24 was released on 21 June 2022.

   The country file has been verified to work with DX4WIN versions 4.0x
   and later. It has also been tested with 3.05, but the information in
   the distance/bearing field in the upper-right of the QSO window will be
   blank.

   The country file contains version information. If you try to log the
   callsign “VERSION” (without the quotes) and hit the tab or space key,
   look at the “QSL Mgr” field in the logbook. This will be the release
   date of the QSL manager database, and will usually match the release
   date of the country file, unless you have imported a different QSL
   manager database.

   To see the actual revision date of the country file, go to Files |
   Databases | Countries. The entry for 1A0, Sov. Military Order of Malta
   should appear. Click on the Callsigns tab, and the callsign “VERSION”
   should appear in the window. The “Start date” is the date the country
   file was created/released.
    1. What’s in the country file?
    2. Installation instructions
    3. Version checking
    4. Merging a new country file with an existing one
    5. Importing QSL manager data
    6. Analyzing country file differences
     __________________________________________________________________

What’s in the Country File?

   The country file unifies many smaller databases that are essential for
   DX4WIN:
     * DXCC entities
     * DXCC prefix to entity mapping
     * DXCC callsign to entity mapping
     * QSL managers
     * QSL manager addresses
     * Islands on the Air (IOTA)
     * US States (latitude/longitude)
     * CT state abbreviation to US state mapping
     * US Counties

   The country file updates are cumulative, not incremental. Each new
   version has all of the previous version’s data (plus/minus any
   corrections), plus the new data for that particular release.
     __________________________________________________________________

Installation instructions

   The new country file can be used “as-is”. First, exit DX4WIN (File |
   Exit); you may be prompted to save your log. Download the ZIP file by
   clicking on the Version History link at the top of this page, then on
   the page that comes up, click on the [download] link of release you
   want. Copy your existing DX4WIN.CTY file (from your DX4W###\SAVE)
   directory to a safe place (like your DX4W###\BACKUP directory, see #3
   in this application note: Analyze DX4WIN Country File Differences).
   Using a utility like WinZip, open the ZIP file you just downloaded, and
   save (extract) the new DX4WIN.CTY file into your DX4W###\SAVE
   directory.

   The ZIP file contains another file, ADIF.PMP. This file is used to map
   QSOs in imported and exported ADIF files to the correct country (this
   includes eQSL and ARRL Logbook of the World). This file must be updated
   any time that DX4WIN.CTY is updated – the files have to be in-sync with
   each other. Assuming the ZIP file is still open from the step above,
   save (extract) the ADIF.PMP file to your DX4W###\IMPORT directory.

   By far, the easiest way to update is to use the DX4WIN Data Updater.
   This updates all DX4WIN data files, including the country file.

   Please see the section below Importing QSL manager data. You may want
   to save your old QSL manager file before installing the new country
   file.
     __________________________________________________________________

Version Checking

   To check if you have installed the most recent country file, type in
   (try to log) the callsign VERSION, then type <TAB> or the spacebar.
   Look at the value in the “Prefix” field in the QSO window. It should
   show the same prefix as the first line in the release notes, i.e.

   Version entity is Morocco, CN

   You can also look at the “QSL Mgr” field in the QSO window.  It should
   have the date of the country file release, i.e.

   08-JAN-2022
     __________________________________________________________________

Merging country files

   If you want to keep a customized country file, but still get the latest
   updates, then merging might be the best solution. Using the new country
   file “as-is” usually results in the fewest problems, however.

   Only the following components of the country file can/will be merged:
     * Exception calls
     * IOTAs
     * QSL managers
     * QSL addresses

   This means entities and prefixes are NOT merged. Therefore, it is
   almost always necessary to merge your old DX4WIN.CTY file into the new
   one, not the other way around. This gets your the prefixes and entities
   that may have been added to the new country file.

   First, exit DX4WIN. Then move your existing DX4WIN.CTY file to a safe
   place, preferably somewhere like DX4W###\BACKUP. If you don’t have a
   BACKUP directory, then create one. See #3 in this application note:
   Analyze DX4WIN Country File Differences.

   Next, download the new country file and extract DX4WIN.CTY to your SAVE
   directory. Start DX4WIN, and ensure you have the new country file by
   trying to log the callsign VERSION and checking the date in the “QSL
   Mgr” field. Close your log using File | Close. If you have made no
   changes (expected) it should not prompt you to save the file.

   To merge, choose Files | Databases | Countries from the main DX4WIN
   menu bar. When the Country Editor window comes up, choose File | Merge
   Other Country File. Navigate to the BACKUP directory, or to wherever
   you moved your old DX4WIN.CTY file. Select the file named DX4WIN.CTY,
   then click on the Open button.

   You will be prompted before merging each of the four components listed
   above. When the merge is complete, check for errors by choosing File |
   Check for Errors in the Country Editor window. When the merging is done
   and there are no errors, choose File | Save Changes and Exit in the
   Country Editor window.

   After the country database is saved, you can re-open your logbook. Do
   whatever test(s) you need to ensure the merge succeeded.
     __________________________________________________________________

Importing QSL manager data

   If you decide you want to use this new file, you may want to import
   your QSL manager data from your old DX4WIN.CTY file into it.

   Using your old country file, go to File | Databases | QSL Managers.
   Then File | Export QSL manager data. Pick a file name and save it. Do
   the same thing for File | Databases | Manager Addresses. Then you can
   exit from the QSL database manager.

   After installing the new DX4WIN.CTY file, go to File | Databases | QSL
   Managers. Then File | Import QSL manager data. Use the same name as you
   saved before. When it asks you if you want to delete the existing data,
   answer Yes. Do the same thing for File | Databases | Manager Addresses.
   Then you can exit from the QSL database manager.

   NOTE: Importing QSL manager data into version 6.02 causes unusual
   dates. This is fixed in version 6.03.

