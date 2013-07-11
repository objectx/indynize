
from __future__ import print_function

import sys, re, os.path, argparse


options = None


def main ():
    global options
    def get_groovy_dir ():
        groovy_dir = os.path.join (os.environ ["HOME"], ".gvm/groovy/current")
        try:
            groovy_dir = os.path.join (os.environ ["GVM_DIR"], "groovy/current")
        except KeyError:
            pass
        return groovy_dir

    parser = argparse.ArgumentParser (description = "")
    parser.add_argument ('dir', nargs = '?', metavar = '<Groovy DIR>', default = get_groovy_dir ())
    parser.add_argument ('-v', '--verbose', action = 'store_true', help = 'Be verbose.')
    parser.add_argument ('-N', '--dry-run', action = 'store_true', help = "Don't modify anything.")

    options = parser.parse_args ()
    # print (options, file = sys.stderr)
    if options.dry_run:
        options.verbose = True

    groovy_dir = options.dir

    libdir = os.path.join (groovy_dir, "lib")
    orig_libdir = os.path.join (groovy_dir, "lib.orig")
    if not os.path.exists (orig_libdir):
        do_rename (libdir, orig_libdir)
    if os.path.exists (libdir):
        do_rmdir (libdir)
    do_mkdir (libdir)
    indy_dir = os.path.join (groovy_dir, "indy")
    # verbose ("libdir = {0}, orig_libdir = {1}, indy_dir = {2}".format (libdir, orig_libdir, indy_dir))
    for f in os.listdir (indy_dir):
        m = re.match ("^(?P<stem>.*)-indy.jar$", f)
        if m:
            stem = m.group ("stem")
            do_link (os.path.join (indy_dir, f), os.path.join (libdir, stem + ".jar"))

    for f in os.listdir (orig_libdir):
        src = os.path.join (orig_libdir, f)
        dst = os.path.join (libdir, f)
        if not os.path.exists (dst):
            do_link (src, dst)
    sys.exit (0)


def verbose (msg):
    if options.verbose:
        print (msg, file = sys.stderr)


def do_rmdir (directory):
    if options.verbose:
        print ("Remove a directory: {0}".format (directory), file = sys.stderr)
    if not options.dry_run:
        os.removedirs (directory)


def do_mkdir (directory):
    if options.verbose:
        print ("Make a directory: {0}".format (directory), file = sys.stderr)
    if not options.dry_run:
        os.mkdir (directory)


def do_rename (src, dst):
    """
    Renames src to dst.

    :param src Source file name:
    :param dst Destination file name:
    """
    if options.verbose:
        print ("Rename {0} to {1}".format (src, dst), file = sys.stderr)
    if not options.dry_run:
        os.rename (src, dst)


def do_link (src, dst):
    """
    Creates hardlink

    :param src Source file name:
    :param dst Destination file name:
    """
    if options.verbose:
        print ("Link {0} as {0}".format (src, dst), file = sys.stderr)
    if not options.dry_run:
        os.link (src, dst)

if __name__ == "__main__":
    main ()

# [END OF FILE]