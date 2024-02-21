import os, sys, itertools
import getopt

# take single file  { lab(key), loc, lat, lon }
# for each label; tss.exe; store last geo value
# use geo value from prevous as anchor for next block
# assumes pre-sorted on label and labels routed


######  Functions ######

# write file for tss.exe to process
def writeTssIn(dat, path):
    with open(path, 'w') as o:
        o.write('lab\tlat\tlon\n')
        for ln in dat:
            o.write('\t'.join(ln) + '\n')

# run tss.exe
def execTss(infile, outfile, anchor="", format="false", image="false", meth="auto"):
    os.system(f'tss.exe -f={infile} -o={outfile} -t={image} -fmt={format} -m={meth} -a={anchor}')  

# get last coords from tss output
def getCoords(path):
    with open(path) as dat:
        lns = [ln.strip() for ln in dat]

    lst = lns[len(lns)-1].split('\t')
    coords = ','.join(lst[1:3])
    return coords

# gen list of rows from file
def readFile(path):
    with open(path) as dat:
        lns = [ln.strip() for ln in dat]

    return lns[1:]   # drop header
    
######  END Functions ######



def main(argv):

    # parse flags
    try:
      opts, args = getopt.getopt(argv,"hi:o:a:")
    except getopt.GetoptError:
      print('daisy.py -i <inputfile> -o <outputfile> -a <starting anchor>')
      sys.exit(2)

    fileIn = 'daisy.dat'
    fileOut = 'out.txt'
    anc = ''

    for opt, arg in opts:
        if opt == '-h':
            print('daisy.py -i <inputfile> -o <outputfile> -a <starting anchor>')
            print('input file: { lab(key), loc, lat, lon }')
            sys.exit()
        elif opt == "-i":
            fileIn = arg
        elif opt == "-o":
            fileOut = arg
        elif opt == "-a":
            anc = arg

    # read dat file
    pwd = os.getcwd()
    try:
        lns = readFile(os.path.join(pwd,fileIn))
    except IOError as e:
        print(f'could not read file.\n {e}')
        sys.exit(2)

    # group data for processing
    groups = []
    for k, g in itertools.groupby(lns, lambda x: x.split('\t')[0]):
        groups.append( [ ln.split('\t')[1:] for ln in g ] )

    # process groups
    out = []
    for i in range(0, len(groups)):

        print(f'\nrouting group {i} with anchor: {anc}')

        # write temp file to process
        tmp_in = 'tmp' + str(i)
        writeTssIn(groups[i], os.path.join(pwd, tmp_in))

        # route temp file
        tmp_o = tmp_in + '_o'
        execTss(tmp_in,tmp_o,anc)

        # get coords for next iter and store
        anc = getCoords(os.path.join(pwd,tmp_o))

        # append result to out list
        tmp_read = readFile(os.path.join(pwd,tmp_o))
        out = out + tmp_read

        # remove temp files
        os.remove(tmp_o)
        os.remove(tmp_in)

    # save out to file
    with open(os.path.join(pwd,'in.txt'), 'w') as o:
        o.write('lab\tlat\tlon\n')
        for ln in out:
            o.write(ln + '\n')

    # build final file and image
    execTss("in.txt",fileOut,"","true","true","none")
    os.remove('in.txt')



if __name__== "__main__":
  main(sys.argv[1:])

