from PIL import Image

def convertTrans(img):
    img = img.convert("RGBA")
    datas = img.getdata()
    new_data = []
    for item in datas:
        if max(item[0], item[1], item[2]) < 60:
            print("Converting to transparent:", item)
            new_data.append((0, 0, 0, 0))
        else:
            new_data.append(item)
    img.putdata(new_data)
    return img
path = "sprites/"
pathnew = "sprites/"
for i in range(0,5):
    for j in range(0,4):

        if i == 0:
            for k in range(0,2):
                img = Image.open(path+str(i)+"/"+str(j)+"-"+str(k)+".png")
                img = convertTrans(img)
                img.save(pathnew+str(i)+"/"+str(j)+"-"+str(k)+".png", "PNG")
        else:
            img = Image.open(path+str(i)+"/"+str(j)+".png")
            img = convertTrans(img)
            img.save(pathnew+str(i)+"/"+str(j)+".png", "PNG")
