/*
bammask : Masking nucleotides in bam files

Copyright © 2022 Institut Pasteur, Paris

Author: Frederic Lemoine

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"log"
	"os"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/sam"

	"github.com/spf13/cobra"
)

var inbam, outbam string
var qual int

// qualityCmd represents the quality command
var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Mask read bases in a bam file, based on base quality",
	Long:  `Mask read nucleotides in a bam file, based on base quality`,
	Run: func(cmd *cobra.Command, args []string) {
		var bamwriter *bam.Writer
		var bamreader *bam.Reader
		var header *sam.Header
		var rec *sam.Record
		var outfile *os.File
		var infile *os.File
		var err error

		// Opening new bam reader
		if inbam == "stdin" || inbam == "-" {
			infile = os.Stdin
		} else {
			if infile, err = os.Open(inbam); err != nil {
				log.Fatal(err)
			}
		}
		if bamreader, err = bam.NewReader(infile, 1); err != nil {
			log.Fatal(err)
		}
		header = bamreader.Header()

		// Opening new bam writer
		if outbam == "stdout" || outbam == "-" {
			outfile = os.Stdout
		} else {
			if outfile, err = os.Create(outbam); err != nil {
				log.Fatal(err)
			}
		}
		if bamwriter, err = bam.NewWriter(outfile, header, 1); err != nil {
			log.Fatal(err)
		}

		// Reading bam file, record by record
		for {
			if rec, err = bamreader.Read(); err != nil {
				if err.Error() != "EOF" {
					log.Print(err)
				}
				break
			}

			// Converting sequence from doublets to one byte per nucleotide
			s := rec.Seq.Expand()
			modified := false
			for i, q := range rec.Qual {
				// If the base at the current index has bad quality, we replace it with a N
				if int(q) < qual {
					modified = true
					s[i] = 'N'
				}
			}
			// If the sequence has been modified, we update it
			if modified {
				rec.Seq = sam.NewSeq(s)
			}

			// We write the record in the output
			if err = bamwriter.Write(rec); err != nil {
				log.Fatal(err)
			}
		}

		bamwriter.Close()
		outfile.Close()
	},
}

func init() {
	rootCmd.AddCommand(qualityCmd)

	qualityCmd.PersistentFlags().StringVarP(&inbam, "input-bam", "i", "stdin", "Input bam file")
	qualityCmd.PersistentFlags().StringVarP(&outbam, "out-bam", "o", "stdout", "Output bam file")
	qualityCmd.PersistentFlags().IntVarP(&qual, "quality", "q", 20, "Quality cutoff below which bases are masked")
}
